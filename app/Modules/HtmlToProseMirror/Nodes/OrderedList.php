<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class OrderedList extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'ol';
    }

    public function data()
    {
        return [
            'type' => 'ordered_list',
            'attrs' => [
                'order' =>
                    $this->DOMNode->getAttribute('start') ?
                    (int) $this->DOMNode->getAttribute('start') :
                    1,
            ],
        ];
    }
}
