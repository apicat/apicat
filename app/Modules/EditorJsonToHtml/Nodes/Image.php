<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class Image extends Node
{
    protected $tagName = 'img';

    public function selfClosing()
    {
        return true;
    }
    
    public function tag()
    {
        return [
            [
                'tag' => 'div',
                'attrs' => [
                    'class' => isset($this->node->attrs->alignment) ? 'image-view image-view--' . $this->node->attrs->alignment : ''
                ]
            ],
            [
                'tag' => 'div',
                'attrs' => [
                    'class' => 'image-view__body'
                ]
            ],
            [
                'tag' => $this->tagName,
                'attrs' => $this->node->attrs,
            ],
        ];
    }
}
