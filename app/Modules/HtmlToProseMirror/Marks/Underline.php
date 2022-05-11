<?php

namespace App\Modules\HtmlToProseMirror\Marks;

class Underline extends Mark
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'u';
    }

    public function data()
    {
        return [
            'type' => 'underline',
        ];
    }
}
