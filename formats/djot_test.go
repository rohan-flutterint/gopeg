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
		{
			Section: "Math",
			Djot:    "Einstein derived $`e=mc^2`.\nPythagoras proved\n$$` x^n + y^n = z^n `",
			Html:    "<p>Einstein derived <span class=\"math inline\">\\(e=mc^2\\)</span>.\nPythagoras proved\n<span class=\"math display\">\\[ x^n + y^n = z^n \\]</span></p>\n",
		},
		{
			Section: "LineBreak",
			Djot:    "This is a soft\nbreak and this is a hard\\\nbreak.",
			Html:    "<p>This is a soft\nbreak and this is a hard<br>\nbreak.</p>\n",
		},
		{
			Section: "LineBreak",
			Djot:    "My reaction is :+1: :smiley:.",
			Html:    "<p>My reaction is 👍 😃.</p>\n",
		},
		{
			Section: "Heading",
			Djot:    "## A level _two_ heading!",
			Html:    "<section id=\"A-level-two-heading\">\n<h2>A level <em>two</em> heading!</h2>\n</section>\n",
		},
		{
			Section: "Heading",
			Djot:    "# A heading that\n# takes up\n# three lines\n\nA paragraph, finally",
			Html:    "<section id=\"A-heading-that-takes-up-three-lines\">\n<h1>A heading that\ntakes up\nthree lines</h1>\n<p>A paragraph, finally</p>\n</section>\n",
		},
		{
			Section: "Heading",
			Djot:    "# A heading that\ntakes up\nthree lines\n\nA paragraph, finally.",
			Html:    "<section id=\"A-heading-that-takes-up-three-lines\">\n<h1>A heading that\ntakes up\nthree lines</h1>\n<p>A paragraph, finally.</p>\n</section>\n",
		},
		//{
		//	Section: "BlockQuote",
		//	Djot:    "> This is a block quote.\n>\n> 1. with a\n> 2. list in it.",
		//	Html:    "<blockquote>\n<p>This is a block quote.</p>\n<ol>\n<li>\nwith a\n</li>\n<li>\nlist in it.\n</li>\n</ol>\n</blockquote>\n",
		//},
		{
			Section: "BlockQuote",
			Djot:    "> This is a block\nquote.",
			Html:    "<blockquote>\n<p>This is a block\nquote.</p>\n</blockquote>\n",
		},
		{
			Section: "BlockQuote",
			Djot:    "> This is a block\n> quote.",
			Html:    "<blockquote>\n<p>This is a block\nquote.</p>\n</blockquote>\n",
		},
		{
			Section: "CodeBlock",
			Djot:    "````\nThis is how you do a code block:\n\n``` ruby\nx = 5 * 6\n```\n````",
			Html:    "<pre><code>This is how you do a code block:\n\n``` ruby\nx = 5 * 6\n```\n</code></pre>\n",
		},
		{
			Section: "CodeBlock",
			Djot:    "``` ruby\nx = 5 * 6\n```",
			Html:    "<pre><code lang=\"ruby\">x = 5 * 6\n</code></pre>\n",
		},
		{
			Section: "BlockQuote",
			Djot:    "> This is a block\n> quote.",
			Html:    "<blockquote>\n<p>This is a block\nquote.</p>\n</blockquote>\n",
		},
		{
			Section: "ThematicBreak",
			Djot:    "Then they went to sleep.\n\n      * * * *\n\nWhen they woke up",
			Html:    "<p>Then they went to sleep.</p>\n<hr>\n<p>When they woke up</p>\n",
		},
		{
			Section: "Div",
			Djot:    "::: warning\nHere is a paragraph.\n\nAnd here is another.\n:::",
			Html:    "<div class=\"warning\">\n<p>Here is a paragraph.</p>\n<p>And here is another.</p>\n</div>\n",
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
