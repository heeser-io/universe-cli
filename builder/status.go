package builder

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

// PrintStatus can be used to print a table of all resources
func (b *Builder) PrintStatus(t *table.Writer) {
	for _, function := range b.cache.Functions {
		(*t).AppendRow(table.Row{b.path, "function", function.Name})
	}
}
