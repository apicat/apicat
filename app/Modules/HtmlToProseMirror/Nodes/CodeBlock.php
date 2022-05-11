<?php

namespace App\Modules\HtmlToProseMirror\Nodes;

class CodeBlock extends Node
{
    public function matching()
    {
        return
            $this->DOMNode->nodeName === 'code' &&
            $this->DOMNode->parentNode->nodeName === 'pre';
    }

    private function getLanguage()
    {
        return preg_replace("/^language-/", "", $this->DOMNode->getAttribute('class'));
    }

    public function data()
    {
        if ($language = $this->getLanguage()) {
            return [
                'type' => 'code_block',
                'attrs' => [
                    'language' => $this->getLanguage(),
                ],
            ];
        }

        return [
            'type' => 'code_block',
        ];
    }
}
