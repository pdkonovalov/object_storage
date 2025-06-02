# Example

```go
package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/pdkonovalov/object_storage"
)

func main() {
	ctx := context.Background()

	cfg := object_storage.Config{
		AccessKey:    "my-access-key",
		SecretKey:    "my-secret-key",
		Region:       "my-region",
		BaseEndpoint: "my-s3-endpoint",
		Bucket:       "my-bucket",
		MetaFilename: "meta.yml",
	}

	s, err := object_storage.New(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Lets file structure of bucket is:
	//
	//   images/file.png
	//   images/meta.yml
	//   file.txt
	//   meta.yml
	//
	// And meta.yml contains:
	//    x: 1
	//    y: 2
	//

	all_bucket, err := s.GetObject(ctx, "")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(all_bucket)

	// Object{
	//    Path:     "",
	//    Meta:     map[string]any{x: 1, y: 2},
	//    Contains: ["images/", "file.txt"],
	//    URL: my-s3-endpoint/my-bucket
	// }

	folder, err := s.GetObject(ctx, "images/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(folder)

	// Object{
	//    Path:     "images/",
	//    Meta:     map[string]any{x: 1, y: 2},
	//    Contains: ["file.png"],
	//    URL: my-s3-endpoint/my-bucket/images/
	// }

	single_file, err := s.GetObject(ctx, "images/file.png")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(single_file)

	// Object{
	//    Path:     "images/file.png",
	//    Meta:     map[string]any{},
	//    Contains: [],
	//    URL: my-s3-endpoint/my-bucket/images/file.png
	// }

	single_file_body, err := s.GetObjectBody(ctx, single_file)
	if err != nil {
		log.Fatal(err)
	}

	defer single_file_body.Close()

	fmt.Println(io.ReadAll(single_file_body))

	//
	// [137 80 78 71 13 10 26 10 0 0 0 13 ...]
	//
}
```
