<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class ApiUrl extends Node
{
    protected $tagName = [
        [
            'tag' => 'div',
            'attrs' => [
                'class' => 'http-url'
            ]
        ]
    ];

    public function text()
    {
        $text = '';
        
        if ($this->node->attrs->url) {
            $text .= '<div class="http-url--url copy_text">' . $this->node->attrs->url . '</div>';
        }
        
        $text .= '<div class="http-url--path"><span class="copy_text">' . $this->node->attrs->path . '</span></div>';
        $text .= '<div class="btn-copy-all copy_text" data-text="' . $this->node->attrs->url . $this->node->attrs->path . '">复制完整URL</div>';

        return $text;
    }
}
