<?php

namespace App\Modules\HtmlToProseMirror\Marks;

class Bold extends Mark
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'strong' || $this->DOMNode->nodeName === 'b';
    }

    public function data()
    {
        return [
            'type' => 'bold',
        ];
    }
}
