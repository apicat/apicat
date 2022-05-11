<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class HardBreak extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'br';
    }

    public function data()
    {
        return [
            'type' => 'hard_break',
        ];
    }
}
