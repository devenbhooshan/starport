package akash

import (
	"bufio"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewScaffold(t *testing.T) {
	t.Run("should scaffold the correct number of files", func(t *testing.T) {
		defer cleanUp()
		c := NewScaffold()
		c.Execute()
		for _, file := range allFiles {
			fileFullPath := strings.TrimSuffix(path.Join("akash", file.fileType, file.fileName), ".tpl")
			_, err := os.Stat(fileFullPath)
			require.NoError(t, err)
		}
	})

	t.Run("scaffolded files should not be empty", func(t *testing.T) {
		defer cleanUp()
		c := NewScaffold()
		c.Execute()

		for _, file := range allFiles {
			fileFullPath := strings.TrimSuffix(path.Join("akash", file.fileType, file.fileName), ".tpl")
			f, err := os.Open(fileFullPath)
			require.NoError(t, err)
			defer f.Close()

			scanner := bufio.NewScanner(f)
			require.True(t, scanner.Scan(), "file is empty", string(fileFullPath))
			require.NotEmpty(t, len(scanner.Text()), "file %s is empty", string(fileFullPath))
		}
	})
}

func cleanUp() {
	os.RemoveAll("akash")
}
