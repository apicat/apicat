<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class OrderedList extends Node
{
    protected $tagName = 'ol';

    public function tag()
    {
        $attrs = [
            'class' => 'ac-ol'
        ];

        if (isset($this->node->attrs->order)) {
            $attrs['start'] = $this->node->attrs->order;
        }

        return [
            [
                'tag' => $this->tagName,
                'attrs' => $attrs,
            ],
        ];
    }
}
