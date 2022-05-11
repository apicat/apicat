<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class Table extends Node
{
    public function matching()
    {
        return
        $this->DOMNode->nodeName === 'tbody' &&
        $this->DOMNode->parentNode->nodeName === 'table';
    }

    public function data()
    {
        return [
            'type' => 'table',
        ];
    }
}
