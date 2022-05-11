<?php

namespace App\Modules\EditorJsonToHtml\Marks;

class Link extends Mark
{
    protected $tagName = 'a';

    public function tag()
    {
        $attrs = [];

        if (isset($this->mark->attrs->openInNewTab) and $this->mark->attrs->openInNewTab) {
            $attrs['target'] = '_blank';
        }

        if (isset($this->mark->attrs->rel)) {
            $attrs['rel'] = $this->mark->attrs->rel;
        }

        $attrs['href'] = $this->mark->attrs->href;

        return [
            [
                'tag' => $this->tagName,
                'attrs' => $attrs,
            ],
        ];
    }
}
