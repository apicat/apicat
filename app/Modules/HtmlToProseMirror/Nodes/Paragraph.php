<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class Paragraph extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'p';
    }

    public function data()
    {
        return [
            'type' => 'paragraph',
        ];
    }
}
