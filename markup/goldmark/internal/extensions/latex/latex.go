package latex

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/util"
	"github.com/yuin/goldmark/text"
)

const LatexTrigger = byte('$')

type LatexParser struct {
}

func (s *LatexParser) Trigger() []byte {
	return []byte{LatexTrigger}
}

func (s *LatexParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	// Finds any LaTeX blocks that start with $ or $$ and returns an AST text node, ensuring
	// that other parsers don't pick up the LaTeX syntax.
	// Source: https://github.com/FurqanSoftware/goldmark-katex/blob/0bf781ec55b4dcac337505863ecccf093e4ee967/parser.go

	buf := block.Source()
	ln, pos := block.Position()

	lstart := pos.Start
	lend := pos.Stop
	line := buf[lstart:lend]

	var start, end, advance int

	if len(line) > 1 && line[1] == LatexTrigger {
		// Block LaTeX
		start = lstart + 2
		offset := 2

	L:
		for x := 0; x < 5; x++ {
			for j := offset; j < len(line); j++ {
				if len(line) > j+1 && line[j] == LatexTrigger && line[j+1] == LatexTrigger {
					end = lstart + j
					advance = 2
					break L
				}
			}
			if lend == len(buf) {
				break
			}
			if end == 0 {
				rest := buf[lend:]
				j := 1
				for j < len(rest) && rest[j] != '\n' {
					j++
				}
				lstart = lend
				lend += j
				line = buf[lstart:lend]
				ln++
				offset = 0
			}
		}

	} else {
		// Inline LaTeX
		start = lstart + 1

		for i := 1; i < len(line); i++ {
			c := line[i]
			if c == '\\' {
				i++
				continue
			}
			if c == LatexTrigger {
				end = lstart + i
				advance = 1
				break
			}
		}
		if end >= len(buf) || buf[end] != LatexTrigger {
			return nil
		}
	}

	if start >= end {
		return nil
	}

	newpos := end + advance
	if newpos < lend {
		block.SetPosition(ln, text.NewSegment(newpos, lend))
	} else {
		block.Advance(newpos)
	}

	return ast.NewTextSegment(text.NewSegment(lstart, newpos))
}

type LatexAsPlainTextExtension struct {
}

func (e *LatexAsPlainTextExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(&LatexParser{}, 0),
	))
}
