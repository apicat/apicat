<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class BulletList extends Node
{
    protected $tagName = [
        [
            'tag' => 'ul',
            'attrs' => [
                'class' => 'ac-ul',
            ]
        ]
    ];
}
