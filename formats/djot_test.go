package formats

import (
	_ "embed"
	"github.com/stretchr/testify/require"
	"testing"
)

type Example struct{ Section, Djot, Html string }

func TestDjot(t *testing.T) {
	for _, tt := range []Example{
		{
			Section: "Link",
			Djot:    "[My link text](http://example.com)",
			Html:    "<p><a href=\"http://example.com\">My link text</a></p>\n",
		},
		{
			Section: "Link",
			Djot:    "[My link text](http://example.com?product_number=234234234234\n234234234234)",
			Html:    "<p><a href=\"http://example.com?product_number=234234234234234234234234\">My link text</a></p>\n",
		},
		{
			Section: "Link",
			Djot:    "[My link text][foo bar]\n\n[foo bar]: http://example.com",
			Html:    "<p><a href=\"http://example.com\">My link text</a></p>\n",
		},
		{
			Section: "Link",
			Djot:    "[foo][bar]",
			Html:    "<p><a>foo</a></p>\n",
		},
		{
			Section: "Link",
			Djot:    "[My link text][]\n\n[My link text]: /url",
			Html:    "<p><a href=\"/url\">My link text</a></p>\n",
		},
		{
			Section: "Image",
			Djot:    "![picture of a cat](cat.jpg)\n\n![picture of a cat][cat]\n\n![cat][]\n\n[cat]: feline.jpg",
			Html:    "<p><img alt=\"picture of a cat\" src=\"cat.jpg\"></p>\n<p><img alt=\"picture of a cat\" src=\"feline.jpg\"></p>\n<p><img alt=\"cat\" src=\"feline.jpg\"></p>\n",
		},
		{
			Section: "AutoLink",
			Djot:    "<https://pandoc.org/lua-filters>\n<me@example.com>",
			Html:    "<p><a href=\"https://pandoc.org/lua-filters\">https://pandoc.org/lua-filters</a>\n<a href=\"mailto:me@example.com\">me@example.com</a></p>\n",
		},
		{
			Section: "Verbatim",
			Djot:    "``Verbatim with a backtick` character``\n`Verbatim with three backticks ``` character`",
			Html:    "<p><code>Verbatim with a backtick` character</code>\n<code>Verbatim with three backticks ``` character</code></p>\n",
		},
		{
			Section: "Verbatim",
			Djot:    "`` `foo` ``",
			Html:    "<p><code>`foo`</code></p>\n",
		},
		{
			Section: "Verbatim",
			Djot:    "`foo bar",
			Html:    "<p><code>foo bar</code></p>\n",
		},
		{
			Section: "Emphasis/strong",
			Djot:    "_emphasized text_\n\n*strong emphasis*",
			Html:    "<p><em>emphasized text</em></p>\n<p><strong>strong emphasis</strong></p>\n",
		},
		{
			Section: "Emphasis/strong",
			Djot:    "_ Not emphasized (spaces). _\n\n___ (not an emphasized `_` character)",
			Html:    "<p>_ Not emphasized (spaces). _</p>\n<p>___ (not an emphasized <code>_</code> character)</p>\n",
		},
		{
			Section: "Emphasis/strong",
			Djot:    "__emphasis inside_ emphasis_",
			Html:    "<p><em><em>emphasis inside</em> emphasis</em></p>\n",
		},
		{
			Section: "Emphasis/strong",
			Djot:    "{_ this is emphasized, despite the spaces! _}",
			Html:    "<p><em> this is emphasized, despite the spaces! </em></p>\n",
		},
		{
			Section: "Highlighted",
			Djot:    "This is {=highlighted text=}.",
			Html:    "<p>This is <mark>highlighted text</mark>.</p>\n",
		},
		{
			Section: "Super/subscript",
			Djot:    "H~2~O and djot^TM^",
			Html:    "<p>H<sub>2</sub>O and djot<sup>TM</sup></p>\n",
		},
		{
			Section: "Super/subscript",
			Djot:    "H{~one two buckle my shoe~}O",
			Html:    "<p>H<sub>one two buckle my shoe</sub>O</p>\n",
		},
		{
			Section: "Insert/delete",
			Djot:    "My boss is {-mean-}{+nice+}.",
			Html:    "<p>My boss is <del>mean</del><ins>nice</ins>.</p>\n",
		},
	} {
		t.Run(tt.Section+"("+tt.Djot+")", func(t *testing.T) {
			html, err := Djot2Html(tt.Djot)
			t.Log(tt.Djot, tt.Html)
			require.Nil(t, err)
			require.Equal(t, tt.Html, html)
		})
	}
}
