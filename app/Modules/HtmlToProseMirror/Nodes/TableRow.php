<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class TableRow extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'tr';
    }

    public function data()
    {
        return [
            'type' => 'tr',
        ];
    }
}
