package main

// import (
// 	"log"

// 	"github.com/JosiahWitt/ensure-cli/internal/packagereader"
// 	"github.com/kr/pretty"
// )

// func main() {
// 	r := packagereader.PackageReader{}
// 	pkgs, err := r.ReadPackages([]*packagereader.PackageDetails{
// 		{
// 			Path:       "github.com/JosiahWitt/ensure-cli/internal/runcmd",
// 			Interfaces: []string{"RunnerIface"},
// 		},
// 		{
// 			Path:       "github.com/JosiahWitt/ensure-cli/internal/fswrite",
// 			Interfaces: []string{"FSWriteIface"},
// 		},
// 		{
// 			Path:       "github.com/JosiahWitt/ensure-cli/internal/exitcleanup",
// 			Interfaces: []string{"ExitCleaner"},
// 		},
// 	})
// 	if err != nil {
// 		log.Fatalln("error", err)
// 	}

// 	pretty.Println(pkgs)
// }
