<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class Heading extends Node
{
    private function getLevel($value)
    {
        preg_match("/^h([1-6])$/", $value, $match);

        return $match[1] ?? null;
    }

    public function matching()
    {
        return (boolean) $this->getLevel($this->DOMNode->nodeName);
    }

    public function data()
    {
        return [
            'type' => 'heading',
            'attrs' => [
                'level' => $this->getLevel($this->DOMNode->nodeName),
            ],
        ];
    }
}
