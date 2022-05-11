<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class HorizontalRule extends Node
{
    protected $tagName = 'hr';

    public function selfClosing()
    {
        return true;
    }
}
