<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class HardBreak extends Node
{
    protected $tagName = 'br';

    public function selfClosing()
    {
        return true;
    }
}
