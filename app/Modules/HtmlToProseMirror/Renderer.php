<?php

namespace App\Modules\HtmlToProseMirror;

use DOMElement;
use DOMDocument;

class Renderer
{
    protected $document;

    protected $storedMarks = [];

    protected $marks = [
        Marks\Bold::class,
        Marks\Code::class,
        Marks\Italic::class,
        Marks\Link::class,
        Marks\Strike::class,
        Marks\Subscript::class,
        Marks\Superscript::class,
        Marks\Underline::class,
    ];

    protected $nodes = [
        Nodes\Blockquote::class,
        Nodes\BulletList::class,
        Nodes\CodeBlock::class,
        Nodes\CodeBlockWrapper::class,
        Nodes\HardBreak::class,
        Nodes\Heading::class,
        Nodes\HorizontalRule::class,
        Nodes\Image::class,
        Nodes\ListItem::class,
        Nodes\OrderedList::class,
        Nodes\Paragraph::class,
        Nodes\Table::class,
        Nodes\TableCell::class,
        Nodes\TableHeader::class,
        Nodes\TableRow::class,
        Nodes\TableWrapper::class,
        Nodes\Text::class,
        Nodes\User::class,
    ];

    public function withMarks($marks = null)
    {
        if (is_array($marks)) {
            $this->marks = $marks;
        }

        return $this;
    }

    public function withNodes($nodes = null)
    {
        if (is_array($nodes)) {
            $this->nodes = $nodes;
        }

        return $this;
    }

    public function addNode($node)
    {
        $this->nodes[] = $node;

        return $this;
    }

    public function addNodes($nodes)
    {
        foreach ($nodes as $node) {
            $this->addNode($node);
        }

        return $this;
    }

    public function addMark($mark)
    {
        $this->marks[] = $mark;

        return $this;
    }

    public function addMarks($marks)
    {
        foreach ($marks as $mark) {
            $this->addMark($mark);
        }

        return $this;
    }

    public function replaceNode($search_node, $replace_node)
    {
        foreach ($this->nodes as $key => $node_class) {
            if ($node_class == $search_node) {
                $this->nodes[$key] = $replace_node;
            }
        }

        return $this;
    }

    public function replaceMark($search_mark, $replace_mark)
    {
        foreach ($this->marks as $key => $mark_class) {
            if ($mark_class == $search_mark) {
                $this->marks[$key] = $replace_mark;
            }
        }

        return $this;
    }

    public function render(string $value): array
    {
        $this->setDocument($value);

        $content = $this->renderChildren(
            $this->getDocumentBody()
        );

        return [
            'type'    => 'doc',
            'content' => $content,
        ];
    }

    private function setDocument(string $value): Renderer
    {
        libxml_use_internal_errors(true);

        $this->document = new DOMDocument;
        $this->document->loadHTML(
            $this->wrapHtmlDocument(
                $this->stripWhitespace($value)
            )
        );

        return $this;
    }

    private function wrapHtmlDocument($value)
    {
        return '<?xml encoding="utf-8" ?>' . $value;
    }

    private function stripWhitespace(string $value): string
    {
        return (new Minify)->process($value);
    }

    private function getDocumentBody(): DOMElement
    {
        return $this->document->getElementsByTagName('body')->item(0);
    }

    private function renderChildren($node): array
    {
        $nodes = [];

        foreach ($node->childNodes as $child) {
            if ($class = $this->getMatchingNode($child)) {
                $item = $class->data();

                if ($item === null) {
                    if ($child->hasChildNodes()) {
                        $nodes = array_merge($nodes, $this->renderChildren($child));
                    }
                    continue;
                }

                if ($child->hasChildNodes()) {
                    $item = array_merge($item, [
                        'content' => $this->renderChildren($child),
                    ]);
                }

                if (count($this->storedMarks)) {
                    $item = array_merge($item, [
                        'marks' => $this->storedMarks,
                    ]);
                }

                if ($class->wrapper) {
                    $item['content'] = [
                        array_merge($class->wrapper, [
                            'content' => @$item['content'] ?: [],
                        ]),
                    ];
                }

                array_push($nodes, $item);
            } elseif ($class = $this->getMatchingMark($child)) {
                array_push($this->storedMarks, $class->data());

                if ($child->hasChildNodes()) {
                    $nodes = array_merge($nodes, $this->renderChildren($child));
                }

                array_pop($this->storedMarks);
            } elseif ($child->hasChildNodes()) {
                $nodes = array_merge($nodes, $this->renderChildren($child));
            }
        }

        return $nodes;
    }

    private function getMatchingNode($item)
    {
        return $this->getMatchingClass($item, $this->nodes);
    }

    private function getMatchingMark($item)
    {
        return $this->getMatchingClass($item, $this->marks);
    }

    private function getMatchingClass($node, $classes)
    {
        foreach ($classes as $class) {
            $instance = new $class($node);

            if ($instance->matching()) {
                return $instance;
            }
        }

        return false;
    }
}
