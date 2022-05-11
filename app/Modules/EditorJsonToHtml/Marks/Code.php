<?php

namespace App\Modules\EditorJsonToHtml\Marks;

class Code extends Mark
{
    protected $tagName = [
        [
            'tag' => 'code',
            'attrs' => [
                'class' => 'code'
            ]
        ]
    ];
}
