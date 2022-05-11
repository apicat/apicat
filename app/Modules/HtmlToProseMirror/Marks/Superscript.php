<?php

namespace App\Modules\HtmlToProseMirror\Marks;

class Superscript extends Mark
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'sup';
    }

    public function data()
    {
        return [
            'type' => 'superscript',
        ];
    }
}
