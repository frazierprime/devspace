package server

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"github.com/devspace-cloud/devspace/sync/remote"
	"github.com/devspace-cloud/devspace/sync/util"
)

var pool = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789%(&)°=?!§ _:$%&/()"

// Generate a random string of A-Z chars with len = l
func random(l int) []byte {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = pool[rand.Intn(len(pool))]
	}
	return bytes
}

type testFile struct {
	Data     []byte
	Children map[string]testFile
}

var fileStructure = testFile{
	Children: map[string]testFile{
		"test.txt": testFile{
			Data: random(10),
		},
		"emptydir": testFile{
			Children: map[string]testFile{},
		},
		"emptydir2": testFile{
			Children: map[string]testFile{},
		},
		"dir1": testFile{
			Children: map[string]testFile{
				"dir1-child": testFile{
					Children: map[string]testFile{
						"test": testFile{
							Data: random(100),
						},
						"test-123": testFile{
							Data: []byte{},
						},
					},
				},
			},
		},
	},
}

var overwriteFileStructure = testFile{
	Children: map[string]testFile{
		"test.txt": testFile{
			Data: random(10),
		},
	},
}

func compareFiles(dir string, file testFile) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	if len(file.Children) != len(files) {
		return fmt.Errorf("dir %s expected %d children, got %d", dir, len(file.Children), len(files))
	}

	// check
	for childName, child := range file.Children {
		found := false
		for _, f := range files {
			if f.Name() == childName {
				if f.IsDir() != (child.Children != nil) {
					return fmt.Errorf("child %s in dir %s: real isDir %v != expected isDir %v", childName, dir, f.IsDir(), child.Children != nil)
				}
				if child.Data != nil {
					data, err := ioutil.ReadFile(filepath.Join(dir, f.Name()))
					if err != nil {
						return err
					}
					if string(data) != string(child.Data) {
						return fmt.Errorf("child %s in dir %s: expected data %s, got %s", childName, dir, string(child.Data), string(data))
					}
				}
				if child.Children != nil {
					err := compareFiles(filepath.Join(dir, childName), child)
					if err != nil {
						return err
					}
				}

				found = true
				break
			}
		}

		if found == false {
			return fmt.Errorf("dir %s: path %s not found", dir, childName)
		}
	}

	return nil
}

func createFiles(dir string, file testFile) error {
	for name, child := range file.Children {
		if child.Children != nil {
			err := os.Mkdir(filepath.Join(dir, name), 0755)
			if err != nil {
				return err
			}

			err = createFiles(filepath.Join(dir, name), child)
			if err != nil {
				return err
			}
		} else {
			err := ioutil.WriteFile(filepath.Join(dir, name), child.Data, 0666)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func TestUpstreamServer(t *testing.T) {
	fromDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}
	toDir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(fromDir)
	defer os.RemoveAll(toDir)

	err = createFiles(fromDir, fileStructure)
	if err != nil {
		t.Fatal(err)
	}

	err = createFiles(toDir, overwriteFileStructure)
	if err != nil {
		t.Fatal(err)
	}

	// Create Upload Tar
	// Open tar
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}

	// Use compression
	gw := gzip.NewWriter(w)
	tarWriter := tar.NewWriter(gw)

	writtenFiles := make(map[string]bool)
	err = recursiveTar(fromDir, "", writtenFiles, tarWriter, false)
	if err != nil {
		t.Fatal(err)
	}

	// Close writer
	tarWriter.Close()
	gw.Close()
	w.Close()

	log.Println("Wrote tar")

	// Upload tar with client
	clientReader, clientWriter := io.Pipe()
	serverReader, serverWriter := io.Pipe()

	go func() {
		err := StartUpstreamServer(toDir, serverReader, clientWriter, false)
		if err != nil {
			t.Fatal(err)
		}
	}()

	conn, err := util.NewClientConnection(clientReader, serverWriter)
	if err != nil {
		t.Fatal(err)
	}

	client := remote.NewUpstreamClient(conn)
	uploadClient, err := client.Upload(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	log.Println("Created server and client")

	// Upload file
	buf := make([]byte, 16*1024)
	for {
		n, err := r.Read(buf)
		if n > 0 {
			err := uploadClient.Send(&remote.Chunk{
				Content: buf[:n],
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		if err == io.EOF {
			_, err := uploadClient.CloseAndRecv()
			if err != nil {
				t.Fatal(err)
			}

			break
		} else if err != nil {
			t.Fatal(err)
		}
	}

	log.Println("Uploaded tar")

	err = compareFiles(toDir, fileStructure)
	if err != nil {
		t.Fatal(err)
	}

	removeClient, err := client.Remove(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	for path := range fileStructure.Children {
		removeClient.Send(&remote.Paths{
			Paths: []string{path, path},
		})
	}

	_, err = removeClient.CloseAndRecv()
	if err != nil {
		t.Fatal(err)
	}

	// Check if toDir is empty
	files, err := ioutil.ReadDir(toDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) > 0 {
		t.Fatalf("Expected empty toDir, but still has %d entries", len(files))
	}
}
