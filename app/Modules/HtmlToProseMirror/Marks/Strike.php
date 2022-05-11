<?php

namespace App\Modules\HtmlToProseMirror\Marks;

class Strike extends Mark
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'strike'
            || $this->DOMNode->nodeName === 's'
            || $this->DOMNode->nodeName === 'del';
    }

    public function data()
    {
        return [
            'type' => 'strike',
        ];
    }
}
