<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class CodeBlock extends Node
{
    public function tag()
    {
        return [
            [
                'tag' => 'div',
                'attrs' => [
                    'class' => 'code-block'
                ]
            ],
            [
                'tag' => 'pre',
                'attrs' => [
                    'data-language' => isset($this->node->attrs->language) ? $this->node->attrs->language : ''
                ]
            ],
            [
                'tag' => 'code'
            ]
        ];
    }

    /**
     * 渲染html开始标签
     *
     * @return string
     */
    public function renderOpeningTag()
    {
        $tags = $this->tag();
        $tags = (array)$tags;

        if (!$tags or count($tags) < 1) {
            return '';
        }

        return implode('', array_map(function($item) {
            if (is_string($item)) {
                return '<' . $item . '>';
            }

            $attrs = '';
            if (isset($item['attrs'])) {
                foreach ($item['attrs'] as $attribute => $value) {
                    $attrs .= ' ' . $attribute . '=' . '"' . $value . '"';
                }
            }

            if ($item['tag'] == 'pre') {
                return '<button class="copy_text">复制</button><' . $item['tag'] . $attrs . '>';
            }

            return '<' . $item['tag'] . $attrs . '>';
        }, $tags));
    }
}
