<?php

namespace App\Repositories\Export;

use Illuminate\Support\Str;
use Illuminate\Support\Facades\File;
use App\Models\Project;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Modules\Mock\Parser\Maker\MockRouter;

/**
 * 导出Postman支持的json文件
 */
class PostmanRepository extends BaseRepository
{
    /**
     * 生成树时要用到的文档字段
     *
     * @var array
     */
    public $docFields = ['title', 'type', 'content'];

    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'postman_export_';

    /**
     * 文件类型
     *
     * @var string
     */
    protected $fileType = 'postman';

    /**
     * 文档默认内容
     *
     * @var array
     */
    public $content;

    /**
     * 请求方法
     *
     * @var array
     */
    protected $requestMethods = [
        '',
        'GET',
        'POST',
        'PUT',
        'PATCH',
        'DELETE',
        'OPTION'
    ];

    /**
     * 数据格式
     *
     * @var array
     */
    protected $bodyModes = [
        '',
        'formdata',
        'urlencoded',
        'raw',
        'binary'
    ];

    /**
     * 构造方法
     *
     * @return void
     */
    public function __construct()
    {
        parent::__construct();
    }

    /**
     * 生成导出文件内容
     *
     * @return boolean
     */
    public function generateContent()
    {
        if (!$project = Project::where('id', $this->projectID)->first()) {
            return $this->fail('导出失败，项目不存在。');
        }

        $this->content = [
            'info' => [
                'name' => $project->name,
                'description' => $project->description ?? '',
                'schema' => 'https://schema.getpostman.com/json/collection/v2.1.0/collection.json'
            ],
            'item' => []
        ];

        if ($this->docID) {
            // 导出单篇文档
            $this->singleDoc();
        } else {
            // 导出整个项目
            $this->project();
        }
    }

    /**
     * 单篇文档导出
     *
     * @return boolean
     */
    protected function singleDoc()
    {
        $doc = ApiDocRepository::getNode($this->docID);
        if (!$doc or !$doc->content) {
            return $this->fail('导出失败，文档不存在。');
        }

        if (!$content = json_decode($doc->content, true)) {
            return $this->fail('导出失败，文档内容有误。');
        }

        $this->content['item'][] = $this->docProcess([
            'title' => $doc->title,
            'content' => $content
        ]);

        $this->makeFile();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = $this->fileType;
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $doc->title) . '_postman';
        $this->cacheContent['projectExport'] = false;

        $this->finish('导出完成');
    }

    /**
     * 整个项目导出
     *
     * @return boolean
     */
    protected function project()
    {
        $tree = $this->makeTree();
        if (!$tree) {
            return $this->fail('导出失败，无法导出一个空的项目。');
        }

        $this->content['item'] = $this->treeProcess($tree);
        
        $this->makeFile();

        if (!$this->result()) {
            return $this->fail('导出失败，请稍后重试。');
        }

        $this->cacheContent['file'] = $this->fileName;
        $this->cacheContent['fileType'] = $this->fileType;
        $this->cacheContent['exportFileName'] = str_replace(' ', '_', $this->content['info']['name']) . '_postman';
        $this->cacheContent['projectExport'] = true;
        $this->finish('导出完成');
    }

    /**
     * 文档树处理
     *
     * @param array $tree 文档树
     * @param array $parent 父级结构
     * @return array
     */
    protected function treeProcess($tree) {
        $result = [];

        foreach ($tree as $v) {
            if ($v['type'] < 1) {
                // 分类
                $result[] = [
                    'name' => $v['title'],
                    'item' => $this->treeProcess($v['sub_nodes'])
                ];
            } else {
                // 文档
                $result[] = $this->docProcess($v);
            }
        }

        return $result;
    }

    /**
     * 将apicat json内容转为postman json内容
     *
     * @param array $doc 文档内容
     * @return array
     */
    protected function docProcess($doc)
    {
        $item = [
            'name' => $doc['title'],
            'request' => [
                'method' => 'GET',
                'header' => []
            ],
            'response' => []
        ];

        if (!$doc['content']) {
            return $item;
        }

        $content = $doc['content'];

        if (count($content['content']) < 1) {
            return $item;
        }

        if ($content['content'][0]['type'] == 'paragraph' and isset($content['content'][0]['content'][0]['text'])) {
            $item['description'] = $content['content'][0]['content'][0]['text'];
        }

        $url = $requestParams = [];
        foreach ($content['content'] as $v) {
            if ($v['type'] == 'http_api_url') {
                $url = $v['attrs']; 
            } elseif ($v['type'] == 'http_api_request_parameter') {
                $requestParams = $v['attrs'];
            }
        }

        if ($url) {         
            $item['request']['method'] = $this->requestMethods[$url['method']];
            $item['request']['body'] = ['mode' => $this->bodyModes[$url['bodyDataType']]];
            $item['request']['url'] = $this->urlProcess($url['url'], $url['path']);
        } else {
            $item['request']['method'] = 'GET';
            $item['request']['url'] = $this->urlProcess('');
        }

        if ($requestParams) {
            $header = $this->baseParamsProcess($requestParams['request_header']['params']);
            if ($header) {
                $item['request']['header'] = $header;
            }

            $query = $this->baseParamsProcess($requestParams['request_query']['params']);
            if (isset($item['request']['url'])) {
                $item['request']['url']['query'] = $query;
                
                $urlQuery = [];
                foreach ($query as $v) {
                    if (strpos($v['key'], '[') === false and isset($v['value'])) {
                        // 过滤掉垃圾内容array object file
                        $urlQuery[] = $v['key'] . '=' . $v['value'];
                    }
                }

                if ($urlQuery) {
                    $item['request']['url']['raw'] .= '?' . implode('&', $urlQuery);
                }
            }

            if (isset($item['request']['body']['mode']) and $item['request']['body']['mode']) {
                if ($item['request']['body']['mode'] == 'raw') {
                    $item['request']['body']['raw'] = json_encode($this->rawParamsProcess($requestParams['request_body']['params']));
                    $item['request']['body']['options'] = ['raw' => ['language' => 'json']];
                } elseif ($item['request']['body']['mode'] == 'binary') {
                    $item['request']['body']['file'] = [
                        'src' => ''
                    ];
                } else {
                    $body = $this->baseParamsProcess($requestParams['request_body']['params']);
                    if ($body) {
                        $item['request']['body'][$item['request']['body']['mode']] = $body;
                    }
                }
            }
        }

        return $item;
    }

    /**
     * url内容处理
     *
     * @param string $raw
     * @param string $path
     * @return array
     */
    protected function urlProcess($raw, $path = '')
    {
        if (!$raw and !$path) {
            return ['raw' => ''];
        }

        $url = [
            'raw' => $raw ? rtrim($raw, '/') . ($path ? '/' . ltrim($path, '/') : '') : '',
        ];

        $protocol = Str::of($raw)->before('://');
        $host = Str::of($raw)->after('://');

        if ($protocol != $raw) {
            $url['protocol'] = (string)$protocol;
        }

        $prePath = Str::of($host)->after('/');
        if ($prePath != $host) {
            $path = rtrim($prePath, '/') . '/' . ltrim($path, '/');
            $host = Str::of($host)->before('/');
        }

        $url['host'] = explode('.', $host);

        if ($path) {
            $url['path'] = explode('/', trim($path, '/'));
        }

        return $url;
    }

    /**
     * 基础参数内容处理
     *
     * @param array $params 参数列表
     * @param array $nameArr 父参数名称数组
     * @param int $parentType 父参数类型
     * @return array
     */
    protected function baseParamsProcess($params, $nameArr = [], $parentType = 5)
    {
        $result = [];

        foreach ($params as $v) {
            $tmp = $nameArr;

            if ($parentType == 4) {
                // 父级参数为array
                $res = preg_match_all('/\[(\w*)\]/', $v['name'], $list);
                if (!$res) {
                    $name = '[]';
                } else {
                    $name = array_pop($list[0]);
                }
            } elseif ($parentType == 5) {
                // 父级参数为object
                $name = $nameArr ? '[' . $v['name'] . ']' : $v['name'];
            } else {
                $name = $v['name'];
            }

            $tmp[] = $name;

            if ($v['type'] == 4 or $v['type'] == 5) {
                $r = $this->baseParamsProcess($v['sub_params'], $tmp, $v['type']);
                if ($r) {
                    $result = array_merge($result, $r);
                }
            } elseif ($v['type'] == 7) {
                $result[] = [
                    'key' => implode('', $tmp),
                    'type' => 'file',
                    'src' => '',
                    'description' => $v['description'] ?? ''
                ];
            } else {
                $result[] = [
                    'key' => implode('', $tmp),
                    'type' => 'text',
                    'value' => $v['default_value'] ? $v['default_value'] : $this->mockValue($v),
                    'description' => $v['description'] ?? ''
                ];
            }
        }

        return $result;
    }

    /**
     * raw参数内容处理
     *
     * @param array $params 参数列表
     * @param int $parentType 父参数类型: 4.array 5.object
     * @return array
     */
    protected function rawParamsProcess($params, $parentType = 5)
    {
        $result = [];

        foreach ($params as $v) {
            switch ($v['type']) {
                case 1:
                    // int
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['default_value'] ? (integer)$v['default_value'] : $this->mockValue($v);
                    } else {
                        $result[] = $v['default_value'] ? (integer)$v['default_value'] : $this->mockValue($v);
                    }
                    break;
                case 2:
                    // float
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['default_value'] ? (float)$v['default_value'] : $this->mockValue($v);
                    } else {
                        $result[] = $v['default_value'] ? (float)$v['default_value'] : $this->mockValue($v);
                    }
                    break;
                case 3:
                    // string
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['default_value'] ? (string)$v['default_value'] : $this->mockValue($v);
                    } else {
                        $result[] = $v['default_value'] ? (string)$v['default_value'] : $this->mockValue($v);
                    }
                    break;
                case 4:
                    // array
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['sub_params'] ? $this->rawParamsProcess($v['sub_params'], 4) : [];
                    } else {
                        $result[] = $v['sub_params'] ? $this->rawParamsProcess($v['sub_params'], 4) : [];
                    }
                    break;
                case 5:
                    // object
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['sub_params'] ? $this->rawParamsProcess($v['sub_params'], 5) : new \stdClass;
                    } else {
                        $result[] = $v['sub_params'] ? $this->rawParamsProcess($v['sub_params'], 5) : new \stdClass;
                    }
                    break;
                case 6:
                    // boolean
                    if ($parentType == 5) {
                        $result[$v['name']] = $v['default_value'] ? 'true' === strtolower($v['default_value']) : $this->mockValue($v);
                    } else {
                        $result[] = $v['default_value'] ? 'true' === strtolower($v['default_value']) : $this->mockValue($v);
                    }
                    break;
                case 7:
                    // file
                    if ($parentType == 5) {
                        $result[$v['name']] = '/file';
                    } else {
                        $result[] = '/file';
                    }
                    break;
            }
        }

        return $result;
    }

    /**
     * mock参数值
     *
     * @param array $param 参数内容
     * @return string|int|float|boolean
     */
    protected function mockValue($param)
    {
        $name = strtolower($param['name']);
        switch ($param['type']) {
            case 1:
                // int
                if (in_array($name, ['mobile', 'idcard', 'zipcode', 'timestamp'])) {
                    return MockRouter::generateData(['type' => 'int', 'mock_type' => $name, 'mock_rule' => '']);
                } else {
                    return MockRouter::generateData(['type' => 'int', 'mock_type' => 'int', 'mock_rule' => '']);
                }
            case 2:
                // float
                return MockRouter::generateData(['type' => 'float', 'mock_type' => 'float', 'mock_rule' => '']);
            case 3:
                // string
                if ($name == 'image' or $name == 'file') {
                    return MockRouter::generateData(['type' => 'string', 'mock_type' => $name . 'url', 'mock_rule' => '']);
                }

                $guessRules = [
                    'mobile', 'phone', 'idcard', 'url', 'domain', 'ip', 'email',
                    'province', 'city', 'zipcode', 'date', 'timestamp'
                ];
                if (in_array($name, $guessRules)) {
                    return MockRouter::generateData(['type' => 'string', 'mock_type' => $name, 'mock_rule' => '']);
                }

                return MockRouter::generateData(['type' => 'string', 'mock_type' => 'string', 'mock_rule' => '']);
            case 6:
                // boolean
                return MockRouter::generateData(['type' => 'boolean', 'mock_type' => 'boolean', 'mock_rule' => '']);
            default:
                return '';
        }
    }

    protected function makeFile()
    {
        if (!File::exists($this->savePath)) {
            File::makeDirectory($this->savePath);
        }
        
        file_put_contents($this->savePath . '/' . $this->fileName, json_encode($this->content, JSON_PRETTY_PRINT));
    }
}
