<?php

namespace App\Modules\HtmlToProseMirror\Marks;

class Mark
{
    public $type = 'mark';

    protected $DOMNode;

    public function __construct($DOMNode)
    {
        $this->DOMNode = $DOMNode;
    }

    public function matching()
    {
        return false;
    }

    public function data()
    {
        return [];
    }
}
