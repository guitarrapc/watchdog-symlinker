package regexast

import (
	"regexp/syntax"
	"strconv"
	"strings"
)

// RegexAst express ast of regex's element.
// Regex character can detect as isRune = true.
type RegexAst struct {
	Len     int
	IsStart bool
	Index   string
	Op      string
	IsRune  bool
	Value   string
}

// ParseRegex parse regexp expression as RE2, and output ast.
// original is regexp.compile : https://github.com/golang/go/blob/faf187fb8e2ca074711ed254c72ffbaed4383c64/src/regexp/regexp.go#L167-L215
func ParseRegex(expr string, mode syntax.Flags) (asts []RegexAst, err error) {
	re, _err := syntax.Parse(expr, mode)
	if _err != nil {
		//fmt.Println(_err)
		return nil, _err
	}
	re = re.Simplify()
	prog, _err := syntax.Compile(re)
	if _err != nil {
		//fmt.Println(_err)
		return nil, _err
	}

	//fmt.Println(prog.String())
	for j := range prog.Inst {
		ast := &RegexAst{}
		i := &prog.Inst[j]
		pc := strconv.Itoa(j)
		if len(pc) < 3 {
			ast.Len = len(pc)
		}
		if j == prog.Start {
			ast.IsStart = true
		}
		// index
		var b strings.Builder
		bw(&b, pc)
		ast.Index = b.String()

		// rune
		ast.Op, ast.IsRune, ast.Value = parseInst(i)
		asts = append(asts, *ast)

		//fmt.Println(ast)
	}
	err = nil
	return
}

func bw(b *strings.Builder, args ...string) {
	for _, s := range args {
		b.WriteString(s)
	}
}

func u32(i uint32) string {
	return strconv.FormatUint(uint64(i), 10)
}

// parse RE2 result with each opcode to get Rune info.
// original is regexp.syntax.dumpInst https://github.com/golang/go/blob/faf187fb8e2ca074711ed254c72ffbaed4383c64/src/regexp/syntax/prog.go#L313-L346
func parseInst(i *syntax.Inst) (op string, isRune bool, value string) {
	op = ""
	isRune = false
	value = ""
	switch i.Op {
	case syntax.InstAlt:
		op = "alt"
		isRune = false
		value = u32(i.Arg)
	case syntax.InstAltMatch:
		op = "altmatch"
		isRune = false
		value = u32(i.Arg)
	case syntax.InstCapture:
		op = "cap"
		isRune = false
		value = u32(i.Arg)
	case syntax.InstEmptyWidth:
		op = "empty"
		isRune = false
		value = u32(i.Arg)
	case syntax.InstMatch:
		op = "match"
		isRune = false
		value = ""
	case syntax.InstFail:
		op = "fail"
		isRune = false
		value = ""
	case syntax.InstNop:
		op = "nop"
		isRune = false
		value = ""
	case syntax.InstRune:
		if i.Rune == nil {
			// shouldn't happen
			op = "rune <nil>"
		}
		op = op + "rune"
		isRune = true
		value = string(i.Rune)
		if syntax.Flags(i.Arg)&syntax.FoldCase != 0 {
			op = op + " /i"
		}
	case syntax.InstRune1:
		op = "rune1"
		isRune = true
		value = string(i.Rune)
	case syntax.InstRuneAny:
		op = "any"
		isRune = false
		value = ""
	case syntax.InstRuneAnyNotNL:
		op = "anynotnl"
		isRune = false
		value = ""
	}

	return
}
