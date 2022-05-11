<?php

namespace App\Modules\EditorJsonToHtml;

class Register
{
    public static $marks = [
        'bold' => Marks\Bold::class,
        'code' => Marks\Code::class,
        'highlight' => Marks\Highlight::class,
        'italic' => Marks\Italic::class,
        'link' => Marks\Link::class,
        'strike' => Marks\Strike::class,
        'subscript' => Marks\Subscript::class,
        'superscript' => Marks\Superscript::class,
        'underline' => Marks\Underline::class,
    ];

    public static $nodes = [
        'blockquote' => Nodes\Blockquote::class,
        'bullet_list' => Nodes\BulletList::class,
        'code_block' => Nodes\CodeBlock::class,
        'hard_break' => Nodes\HardBreak::class,
        'heading' => Nodes\Heading::class,
        'horizontal_rule' => Nodes\HorizontalRule::class,
        'image' => Nodes\Image::class,
        'list_item' => Nodes\ListItem::class,
        'ordered_list' => Nodes\OrderedList::class,
        'paragraph' => Nodes\Paragraph::class,
        'table' => Nodes\Table::class,
        'td' => Nodes\TableCell::class,
        'th' => Nodes\TableHeader::class,
        'tr' => Nodes\TableRow::class,
        'http_api_url' => Nodes\HttpApiUrl::class,
        'api_url' => Nodes\ApiUrl::class,
        'api_parameter' => Nodes\ApiParameter::class,
        'http_api_request_parameter' => Nodes\HttpApiRequestParameter::class,
        'http_api_response_parameter' => Nodes\HttpApiResponseParameter::class,
        'http_status_code' => Nodes\HttpStatusCode::class,
    ];

    /**
     * 实例化mark
     *
     * @param object $mark 节点对象
     * @return Marks\Mark
     */
    public static function mark($mark)
    {
        if (!isset($mark->type) or !isset(self::$marks[$mark->type])) {
            return;
        }

        return new self::$marks[$mark->type]($mark);
    }

    /**
     * 实例化Node
     *
     * @param object $node 节点对象
     * @param boolean $print 是否是打印格式
     * @return Nodes\Node
     */
    public static function node($node, $print)
    {
        if (!isset($node->type) or !isset(self::$nodes[$node->type])) {
            return;
        }

        return new self::$nodes[$node->type]($node, $print);
    }
}
