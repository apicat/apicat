<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class Image extends Node
{
    public function matching()
    {
        return $this->DOMNode->nodeName === 'img';
    }

    public function data()
    {
        return [
            'type' => 'image',
            'attrs' => [
                'alt' => $this->DOMNode->hasAttribute('alt') ? $this->DOMNode->getAttribute('alt') : null,
                'src' => $this->DOMNode->hasAttribute('src') ? $this->DOMNode->getAttribute('src') : null,
                'title' => $this->DOMNode->hasAttribute('title') ? $this->DOMNode->getAttribute('title') : null,
            ],
        ];
    }
}
