<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class HttpStatusCode extends Node
{
    protected $tagName = [
        [
            'tag' => 'div',
            'attrs' => [
                'class' => 'http-code'
            ]
        ]
    ];

    public function text()
    {
        $codeTagColor = ['', '#BEBEBE', '#66BE74', '#51B9C3', '#F1924E', '#DF4545'];
        $colorIndex = floor($this->node->attrs->code / 100);

        $htmls = '<span class="intro">' . $this->node->attrs->intro . '</span>';
        $htmls .= '<span class="code" style="background: ' . $codeTagColor[$colorIndex] . ';" data-tippy-content="' . $this->node->attrs->codeDesc . '">' . $this->node->attrs->code . '</span>';
        return $htmls;
    }
}
