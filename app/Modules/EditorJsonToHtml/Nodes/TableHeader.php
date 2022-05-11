<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class TableHeader extends TableCell
{
    protected $tagName  = 'th';

    public function tag()
    {
        $attrs = [];
        if (isset($this->node->attrs)) {
            if (isset($this->node->attrs->colspan)) {
                $attrs['colspan'] = $this->node->attrs->colspan;
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
