<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class HttpApiUrl extends Node
{
    protected $tagName = [
        [
            'tag' => 'div',
            'attrs' => [
                'class' => 'http-url'
            ]
        ]
    ];

    protected $methodTagColor = [
        '',
        '#66BE74',
        '#4894FF',
        '#51B9C3',
        '#F1924E',
        '#DF4545',
        '#A973DF'
    ];

    protected $methodTagName = [
        '',
        'GET',
        'POST',
        'PUT',
        'PATCH',
        'DELETE',
        'OPTION'
    ];

    protected $bodyDataType = [
        'none',
        'form-data',
        'x-www-form-urlencoded',
        'raw',
        'binary'
    ];

    public function text()
    {
        if ($this->print) {
            return $this->printText();
        }

        $text = '<div class="http-url--method" style="background: ' . $this->methodTagColor[$this->node->attrs->method] . ';">' . $this->methodTagName[$this->node->attrs->method] . '</div>';

        if (isset($this->node->attrs->bodyDataType) and $this->node->attrs->method != 1) {
            $text .= '<div class="http-url--type">' . $this->bodyDataType[$this->node->attrs->bodyDataType] . '</div>';
        }

        if ($this->node->attrs->url) {
            $text .= '<div class="http-url--url copy_text">' . $this->node->attrs->url . '</div>';
        }
        
        $text .= '<div class="http-url--path"><span class="copy_text">' . $this->node->attrs->path . '</span></div>';
        $text .= '<div class="btn-copy-all copy_text" data-text="' . $this->node->attrs->url . $this->node->attrs->path . '">复制完整URL</div>';

        return $text;
    }

    public function printText()
    {
        return '';
    }
}
