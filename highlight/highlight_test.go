package highlight

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPython(t *testing.T) {
	highlighted, err := Highlight(`def square(value): 
    result = value**2 # dummy comment
    return result`, PythonTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="keyword">def</span> <span class="function">square</span>(<span class="identifier">value</span>): 
    <span class="identifier">result</span> = <span class="identifier">value</span>**2 <span class="comment"># dummy comment</span>
    <span class="keyword">return</span> <span class="identifier">result</span>`, highlighted)
}

func TestC(t *testing.T) {
	highlighted, err := Highlight(`typedef struct {
    unsigned value;       /**comment */
} parameters;`, CTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="keyword">typedef</span> <span class="keyword">struct</span> {
    <span class="keyword">unsigned</span> <span class="identifier">value</span>;       <span class="comment">/**comment */</span>
} <span class="identifier">parameters</span>;`, highlighted)
}

func TestRust(t *testing.T) {
	highlighted, err := Highlight(`fn main() { println!("Hello, world!"); }`, RustTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="keyword">fn</span> <span class="function">main</span>() { <span class="identifier">println</span>!(<span class="string">"Hello, world!"</span>); }`, highlighted)
}

func TestShell(t *testing.T) {
	highlighted, err := Highlight(`$> echo hi
123`, ShellTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="command">$> echo hi</span>
123`, highlighted)
}

func TestGo(t *testing.T) {
	highlighted, err := Highlight(`type E struct{ Desc string }

func (e *E) Error() string { return e.Desc }
func Api() *E              { return nil }

func TestApi(t *testing.T) {
  var err error
  err = Api()
  require.True(t, err == nil)
}`, GoTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="keyword">type</span> <span class="identifier">E</span> <span class="keyword">struct</span>{ <span class="identifier">Desc</span> <span class="identifier">string</span> }

<span class="keyword">func</span> (<span class="identifier">e</span> *<span class="identifier">E</span>) <span class="function">Error</span>() <span class="identifier">string</span> { <span class="keyword">return</span> <span class="identifier">e</span>.<span class="identifier">Desc</span> }
<span class="keyword">func</span> <span class="function">Api</span>() *<span class="identifier">E</span>              { <span class="keyword">return</span> <span class="identifier">nil</span> }

<span class="keyword">func</span> <span class="function">TestApi</span>(<span class="identifier">t</span> *<span class="identifier">testing</span>.<span class="identifier">T</span>) {
  <span class="keyword">var</span> <span class="identifier">err</span> <span class="identifier">error</span>
  <span class="identifier">err</span> = <span class="function">Api</span>()
  <span class="identifier">require</span>.<span class="function">True</span>(<span class="identifier">t</span>, <span class="identifier">err</span> == <span class="identifier">nil</span>)
}`, highlighted)
}

func TestAsm(t *testing.T) {
	highlighted, err := Highlight(`TEXT     main.Check(SB), NOSPLIT|NOFRAME|ABIInternal, $0-0
FUNCDATA $0, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
FUNCDATA $1, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
XORL     AX, AX
RET`, AsmTokenizerRules)
	require.Nil(t, err)
	require.Equal(t, `<span class="keyword">TEXT</span>     main.Check(SB), NOSPLIT|NOFRAME|ABIInternal, <span class="number">$0</span><span class="number">-0</span>
<span class="keyword">FUNCDATA</span> <span class="number">$0</span>, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
<span class="keyword">FUNCDATA</span> <span class="number">$1</span>, gclocals·g2BeySu+wFnoycgXfElmcg==(SB)
<span class="keyword">XORL</span>     AX, AX
<span class="keyword">RET</span>`, highlighted)
}
