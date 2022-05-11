<?php

namespace App\Modules\EditorJsonToHtml\Marks;

class Highlight extends Mark
{
    protected $tagName = 'mark';

    public function tag()
    {
        $attrs = [];
        $style = '';

        if (isset($this->mark->attrs->bgColor) and $this->mark->attrs->bgColor) {
            $style .= 'background-color: ' . $this->mark->attrs->bgColor .';';
        }

        if (isset($this->mark->attrs->fontColor) and $this->mark->attrs->fontColor) {
            $style .= 'color: ' . $this->mark->attrs->fontColor .';';
        }

        if ($style) {
            $attrs['style'] = $style;
        }

        return [
            [
                'tag' => $this->tagName,
                'attrs' => $attrs,
            ],
        ];
    }
}
