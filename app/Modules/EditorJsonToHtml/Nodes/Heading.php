<?php

namespace App\Modules\EditorJsonToHtml\Nodes;

class Heading extends Node
{
    public function tag()
    {
        return "h{$this->node->attrs->level}";
    }
}
