package gen

import (
	"github.com/stretchr/testify/require"
	"go/parser"
	"go/token"
	"golang.org/x/mod/modfile"
	"io"
	"os"
	"testing"
)

func TestCreateEmptyProject(t *testing.T) {
	t.Run("test create empty project", func(t *testing.T) {
		err := createEmptyProject("github.com/test/test", "test")
		require.NoError(t, err)

		var content []byte
		{
			f, err := os.Open("test/go.mod")
			require.NoError(t, err)
			content, err = io.ReadAll(f)
			require.NoError(t, err)
			f.Close()
		}

		{
			_, err := modfile.Parse("test/go.mod", content, nil)
			require.NoError(t, err)
		}
	})

	t.Cleanup(func() {
		os.RemoveAll("test")
	})
}

func TestCreateMainFile(t *testing.T) {
	t.Run("test create main file", func(t *testing.T) {
		err := createGoFile("./main.go", false)
		require.NoError(t, err)

		{
			fset := token.NewFileSet()
			_, err := parser.ParseFile(fset, "main.go", nil, parser.ParseComments)
			require.NoError(t, err)
			//ast.Print(fset, f)
		}
	})

	t.Cleanup(func() {
		os.RemoveAll("main.go")
	})
}
