<?php

namespace App\Modules\Editor;

/**
 * 编辑器文档内容
 */
class Content
{
    /**
     * 内容
     *
     * @var array
     */
    public $content;

    /**
     * http 响应状态码
     *
     * @var array
     */
    protected $httpCode = [
        100 => 'Continue',
        101 => 'Switching Protocols',
        102 => 'Processing',
        103 => 'Early Hints',
        200 => 'OK',
        201 => 'Created',
        202 => 'Accepted',
        203 => 'Non-Authoritative Information',
        204 => 'No Content',
        205 => 'Reset Content',
        206 => 'Partial Content',
        207 => 'Multi-Status',
        208 => 'Already Reported',
        226 => 'IM Used',
        300 => 'Multiple Choices',
        301 => 'Moved Permanently',
        302 => 'Found',
        303 => 'See Other',
        304 => 'Not Modified',
        305 => 'Use Proxy',
        306 => 'Switch Proxy',
        307 => 'Temporary Redirect',
        308 => 'Permanent Redirect',
        400 => 'Bad Request',
        401 => 'Unauthorized',
        402 => 'Payment Required',
        403 => 'Forbidden',
        404 => 'Not Found',
        405 => 'Method Not Allowed',
        406 => 'Not Acceptable',
        407 => 'Proxy Authentication Required',
        408 => 'Request Timeout',
        409 => 'Conflict',
        410 => 'Gone',
        411 => 'Length Required',
        412 => 'Precondition Failed',
        413 => 'Request Entity Too Large',
        414 => 'Request-URI Too Long',
        415 => 'Unsupported Media Type',
        416 => 'Requested Range Not Satisfiable',
        417 => 'Expectation Failed',
        418 => "I'm a teapot",
        420 => 'Enhance Your Calm',
        421 => 'Misdirected Request',
        422 => 'Unprocessable Entity',
        423 => 'Locked',
        424 => 'Failed Dependency',
        425 => 'Too Early',
        426 => 'Upgrade Required',
        428 => 'Precondition Required',
        429 => 'Too Many Requests',
        431 => 'Request Header Fields Too Large',
        444 => 'No Response',
        450 => 'Blocked by Windows Parental Controls',
        451 => 'Unavailable For Legal Reasons',
        494 => 'Request Header Too Large',
        500 => 'Internal Server Error',
        501 => 'Not Implemented',
        502 => 'Bad Gateway',
        503 => 'Service Unavailable',
        504 => 'Gateway Timeout',
        505 => 'HTTP Version Not Supported',
        506 => 'Variant Also Negotiates',
        507 => 'Insufficient Storage',
        508 => 'Loop Detected',
        510 => 'Not Extended',
        511 => 'Network Authentication Required'
    ];

    public function __construct()
    {
        $this->init();
    }

    /**
     * 初始化内容
     *
     * @return void
     */
    public function init()
    {
        $this->content = [
            'type' => 'doc',
            'content' => []
        ];
    }

    /**
     * 获取文档Json内容
     *
     * @return string
     */
    public function contentJson()
    {
        return json_encode($this->content);
    }

    /**
     * 获取原始内容
     *
     * @return array
     */
    public function contentOriginal()
    {
        return $this->content;
    }

    /**
     * 添加标题
     *
     * @param string $text 文本内容
     * @param integer $level 标题等级: 1.H1 2.H2 3.H3 4.H4 5.H5 6.H6
     * @return void
     */
    public function addHeading($text, $level)
    {
        $this->addContnet([
            'type' => 'heading',
            'attrs' => ['level' => $level],
            'content' => [
                [
                    'type' => 'text',
                    'text' => $text
                ]
            ]
        ]);
    }

    /**
     * 添加段落
     *
     * @param string $text 文本内容
     * @return void
     */
    public function addParagraph($text)
    {
        $this->addContnet([
            'type' => 'paragraph',
            'content' => [
                [
                    'type' => 'text',
                    'text' => $text
                ]
            ]
        ]);
    }

    /**
     * 添加HTTP API URL组件及内容
     *
     * @param string $url 请求地址
     * @param string $path 请求路径
     * @param integer $method 请求方法: 1.GET 2.POST 3.PUT 4.PATCH 5.DELETE 6.OPTION
     * @param integer $bodyDataType 请求体参数类型: 1.none 2.form-data 3.x-www-form-urlencoded 4.raw 5.binary
     * @return void
     */
    public function addHttpApiUrl($url, $path, $method = 1, $bodyDataType = 1)
    {
        $this->addContnet([
            'type' => 'http_api_url',
            'attrs' => [
                'url' => $url,
                'path' => $path,
                'method' => $method,
                'bodyDataType' => $bodyDataType
            ]
        ]);
    }

    /**
     * 添加HTTP API请求参数
     *
     * @param string $headerTitle 请求头名称
     * @param array $header 请求头参数
     * @param string $bodyTitle 请求体名称
     * @param array $body 请求体参数
     * @param string $queryTitle URL请求参数名称
     * @param array $query URL请求参数
     * @return void
     */
    public function addHttpApiRequestParams($headerTitle, $header, $bodyTitle, $body, $queryTitle, $query)
    {
        $this->addContnet([
            'type' => 'http_api_request_parameter',
            'attrs' => [
                'request_header' => [
                    'title' => $headerTitle,
                    'params' => Helper\RequestParams::filter($header)
                ],
                'request_body' => [
                    'title' => $bodyTitle,
                    'params' => Helper\RequestParams::filter($body)
                ],
                'request_query' => [
                    'title' => $queryTitle,
                    'params' => Helper\RequestParams::filter($query)
                ]
            ]
        ]);
    }

    /**
     * 添加HTTP API返回参数
     *
     * @param string $headerTitle 返回头名称
     * @param array $header 返回头参数
     * @param string $bodyTitle 返回体名称
     * @param array $body 返回体参数
     * @return void
     */
    public function addHttpApiResponseParams($headerTitle, $header, $bodyTitle, $body)
    {
        $this->addContnet([
            'type' => 'http_api_response_parameter',
            'attrs' => [
                'response_header' => [
                    'title' => $headerTitle,
                    'params' => Helper\ResponseParams::filter($header)
                ],
                'response_body' => [
                    'title' => $bodyTitle,
                    'params' => Helper\ResponseParams::filter($body)
                ]
            ]
        ]);
    }

    /**
     * 添加HTTP API 返回状态码
     *
     * @param int $code 状态码
     * @return void
     */
    public function addHttpApiResponseStatusCode($code)
    {
        if (!isset($this->httpCode[$code])) {
            return;
        }

        $this->addContnet([
            'type' => 'http_status_code',
            'attrs' => [
                'intro' => 'Response Status Code:',
                'code' => $code,
                'codeDesc' => $this->httpCode[$code]
            ]
        ]);
    }

    /**
     * 添加代码块
     *
     * @param string $language 代码语言
     * @param string $content 代码内容
     * @return void
     */
    public function addCodeBlock($language, $content)
    {
        $this->addContnet([
            'type' => 'code_block',
            'attrs' => [
                'language' => $language
            ],
            'content' => [
                [
                    'type' => 'text',
                    'text' => $content
                ]
            ]
        ]);
    }

    /**
     * 添加内容
     *
     * @param array $content block内容
     * @return void
     */
    protected function addContnet($content)
    {
        $this->content['content'][] = $content;
    }
}