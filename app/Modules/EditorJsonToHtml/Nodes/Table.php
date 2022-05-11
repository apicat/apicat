<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class Table extends Node
{
    protected $tagName = [
        [
            'tag' => 'div',
            'attrs' => [
                'class' => 'scrollable-wrapper original-table'
            ]
        ],
        [
            'tag' => 'div',
            'attrs' => [
                'class' => 'scrollable'
            ]
        ],
        [
            'tag' => 'table',
            'attrs' => [
                'class' => 'rme-table'
            ]
        ],
        [
            'tag' => 'tbody'
        ]
    ];
}
