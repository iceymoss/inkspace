package config

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func readDeploymentYAML(t *testing.T, path string, target any) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("../..", path))
	if err != nil {
		t.Fatal(err)
	}
	if err := yaml.Unmarshal(body, target); err != nil {
		t.Fatal(err)
	}
}

func TestRenderedKubernetesCanonicalReferences(t *testing.T) {
	if _, err := exec.LookPath("kubectl"); err != nil {
		t.Skip("kubectl not installed; source manifests still checked")
	}
	for _, app := range []string{"blog-api", "admin-api", "scheduler-task"} {
		cmd := exec.Command("kubectl", "kustomize", "../../deploy/apps/backend/"+app+"/overlays/development")
		output, err := cmd.Output()
		if err != nil {
			t.Fatal(err)
		}
		decoder := yaml.NewDecoder(strings.NewReader(string(output)))
		var configName, reference string
		for {
			var doc map[string]any
			if err := decoder.Decode(&doc); err != nil {
				if err == io.EOF {
					break
				}
				t.Fatal(err)
			}
			metadata, ok := doc["metadata"].(map[string]any)
			if !ok {
				t.Fatal("missing metadata")
			}
			if metadata["namespace"] != "inkspace-backend" {
				t.Fatalf("%s resource in wrong namespace: %v", app, metadata)
			}
			if doc["kind"] == "ConfigMap" {
				configName, _ = metadata["name"].(string)
				data, _ := doc["data"].(map[string]any)
				if data["PUBLIC_SITE_URL"] != "http://dev.blog.test.com" {
					t.Fatal(data)
				}
			}
			if doc["kind"] == "Deployment" {
				// Decode only this known Deployment's env shape rather than assume
				// YAML map nesting/type assertions throughout the test.
				encoded, err := yaml.Marshal(doc)
				if err != nil {
					t.Fatal(err)
				}
				var node yaml.Node
				if err := yaml.Unmarshal(encoded, &node); err != nil {
					t.Fatal(err)
				}
				var findReference func(*yaml.Node)
				findReference = func(n *yaml.Node) {
					if n.Kind == yaml.MappingNode {
						for i := 0; i < len(n.Content); i += 2 {
							if n.Content[i].Value == "configMapKeyRef" {
								var ref struct{ Name string }
								if err := n.Content[i+1].Decode(&ref); err != nil {
									t.Fatal(err)
								}
								reference = ref.Name
							}
						}
					}
					for _, child := range n.Content {
						findReference(child)
					}
				}
				findReference(&node)
			}
		}
		if configName == "" || reference != configName || configName == app+"-public-site" {
			t.Fatalf("%s generated config=%q deployment ref=%q", app, configName, reference)
		}
	}
}

func TestComposeEveryBinaryRequiresCanonicalOrigin(t *testing.T) {
	for _, path := range []string{"docker-compose.yml", "docker-compose.external-db.yml"} {
		t.Run(path, func(t *testing.T) {
			var manifest struct {
				Services map[string]struct{ Environment yaml.Node }
			}
			readDeploymentYAML(t, path, &manifest)
			for _, name := range []string{"backend-1", "backend-2", "backend-3", "admin-backend", "scheduler"} {
				var environment []string
				service := manifest.Services[name]
				if err := service.Environment.Decode(&environment); err != nil {
					t.Fatal(err)
				}
				found := false
				for _, variable := range environment {
					if strings.HasPrefix(variable, "PUBLIC_SITE_URL=${PUBLIC_SITE_URL:?") {
						found = true
					}
				}
				if !found {
					t.Errorf("%s does not require and forward PUBLIC_SITE_URL", name)
				}
			}
		})
	}
}

func TestKubernetesEveryBinaryUsesEnvironmentSpecificCanonicalConfig(t *testing.T) {
	for _, app := range []string{"blog-api", "admin-api", "scheduler-task"} {
		t.Run(app, func(t *testing.T) {
			path := "deploy/apps/backend/" + app
			var deployment struct {
				Spec struct {
					Template struct {
						Spec struct {
							Containers []struct {
								Env []struct {
									Name      string
									Value     string
									ValueFrom struct {
										ConfigMapKeyRef struct {
											Name, Key string
											Optional  bool
										} `yaml:"configMapKeyRef"`
									} `yaml:"valueFrom"`
								}
							}
						}
					}
				}
			}
			readDeploymentYAML(t, path+"/base/deployment.yaml", &deployment)
			if len(deployment.Spec.Template.Spec.Containers) != 1 {
				t.Fatal("expected application container")
			}
			found := false
			for _, env := range deployment.Spec.Template.Spec.Containers[0].Env {
				if env.Name == "PUBLIC_SITE_URL" {
					ref := env.ValueFrom.ConfigMapKeyRef
					found = env.Value == "" && ref.Name == app+"-public-site" && ref.Key == "PUBLIC_SITE_URL" && !ref.Optional
				}
			}
			if !found {
				t.Fatal("base must require environment-owned ConfigMap, never a baked development origin")
			}
			var overlay struct {
				ConfigMapGenerator []struct {
					Name     string
					Literals []string
				} `yaml:"configMapGenerator"`
			}
			readDeploymentYAML(t, path+"/overlays/development/kustomization.yaml", &overlay)
			if len(overlay.ConfigMapGenerator) != 1 || overlay.ConfigMapGenerator[0].Name != app+"-public-site" || len(overlay.ConfigMapGenerator[0].Literals) != 1 || overlay.ConfigMapGenerator[0].Literals[0] != "PUBLIC_SITE_URL=http://dev.blog.test.com" {
				t.Fatal("development must explicitly supply the same BLOG origin for every binary")
			}
		})
	}
}
