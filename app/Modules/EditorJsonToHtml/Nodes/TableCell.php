<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class TableCell extends Node
{
    protected $tagName = 'td';

    public function tag()
    {
        $attrs = [];
        if (isset($this->node->attrs)) {
            if (isset($this->node->attrs->colspan)) {
                $attrs['colspan'] = $this->node->attrs->colspan;
            }

            if (isset($this->node->attrs->colwidth)) {
                if ($widths = $this->node->attrs->colwidth) {
                    if (count($widths) === $attrs['colspan']) {
                        $attrs['data-colwidth'] = implode(',', $widths);
                    }
                }
            }

            if (isset($this->node->attrs->rowspan)) {
                $attrs['rowspan'] = $this->node->attrs->rowspan;
            }

            if (isset($this->node->attrs->alignment)) {
                $attrs['style'] = 'text-align: ' . $this->node->attrs->alignment;
            }
        }

        return [
            [
                'tag' => $this->tagName,
                'attrs' => $attrs,
            ]
        ];
    }
}
