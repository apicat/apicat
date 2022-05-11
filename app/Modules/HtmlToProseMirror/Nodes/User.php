<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class User extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'user-mention';
    }

    public function data()
    {
        return [
            'type' => 'user',
            'attrs' => [
                'id' => $this->DOMNode->getAttribute('data-id'),
            ],
        ];
    }
}
