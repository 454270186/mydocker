/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/454270186/mydocker/cmd"
	_ "github.com/454270186/mydocker/nsenter"
)

func main() {
	cmd.Execute()
}
