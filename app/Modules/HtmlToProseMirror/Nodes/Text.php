<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class Text extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === '#text';
    }

    public function data()
    {
        $text = ltrim($this->DOMNode->nodeValue, "\n");

        if ($text === '') {
            return null;
        }

        return [
            'type' => 'text',
            'text' => $text,
        ];
    }
}
