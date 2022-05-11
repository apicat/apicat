<?php

namespace App\Repositories\Import;

use Illuminate\Support\Facades\File;
use App\Repositories\Project\ApiDocRepository;
use App\Modules\Editor\Content;
use App\Modules\Editor\Helper\ParamMaker;

/**
 * PostMan Json文件导入
 */
class PostmanRepository extends BaseRepository
{
    /**
     * 缓存key前缀
     *
     * @var string
     */
    protected $cachePrefix = 'postman_import_';

    /**
     * apicat文档内容对象
     *
     * @var Content
     */
    protected $apicatContentObj;

    /**
     * 参数处理对象
     *
     * @var ParamMaker
     */
    protected $paramMakerObj;

    /**
     * 请求方法
     *
     * @var array
     */
    protected $requestMethods = [
        'GET' => 1,
        'POST' => 2,
        'PUT' => 3,
        'PATCH' => 4,
        'DELETE' => 5,
        'OPTION' => 6
    ];

    /**
     * 数据格式
     *
     * @var array
     */
    protected $bodyModes = [
        'none' => 0,
        'formdata' => 1,
        'urlencoded' => 2,
        'raw' => 3,
        'file' => 4
    ];

    public function __construct()
    {
        parent::__construct();
        $this->apicatContentObj = new Content;
        $this->paramMakerObj = new ParamMaker;
    }

    /**
     * 读取文件内容
     *
     * @return void
     */
    protected function readFile()
    {
        if (!File::exists($this->filePath)) {
            return $this->fail('导入失败，请重新导入。');
        }
        
        if (!$content = file_get_contents($this->filePath)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，请重新导入。');
        }

        if (!$content = trim($content)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        if (!$content = json_decode($content, true)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        if (!is_array($content)) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        if (!isset($content['info'], $content['item'])) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        if (!is_array($content['item']) or !$content['item']) {
            File::delete($this->filePath);
            return $this->fail('导入失败，文件内容有误。');
        }

        $this->import($content['item'], $this->parentID);

        $this->finish('导入完成');
        File::delete($this->filePath);
    }

    /**
     * 导入文档
     *
     * @param array $content 文档内容
     * @param integer $parentID 父级id
     * @return void
     */
    protected function import($docs, $parentID = 0)
    {
        foreach ($docs as $doc) {
            if (!isset($doc['name'])) {
                continue;
            }

            if (isset($doc['item'])) {
                // 分类
                $record = ApiDocRepository::addDirToFoot($this->projectID, $doc['name'], $parentID, $this->userID);

                if ($doc['item']) {
                    $this->import($doc['item'], $record->id);
                }
            } else {
                // 文档
                if (!isset($doc['request'])) {
                    continue;
                }

                $content = $this->process($doc);
                ApiDocRepository::addDoc($this->projectID, $parentID, $doc['name'], $content, $this->userID);
            }
        }
    }

    /**
     * 处理文档内容
     *
     * @param array $content 原文内容
     * @return string
     */
    protected function process($content)
    {
        $this->apicatContentObj->init();

        if (isset($content['request']) and $content['request']) {
            $this->requestProcess($content['request']);
        }
        
        if (isset($content['response']) and is_array($content['response']) and $content['response']) {
            $this->responseProcess($content['response'][0]);
        }

        return $this->apicatContentObj->contentJson();
    }

    /**
     * 请求内容处理
     *
     * @param array $request 请求内容
     * @return void
     */
    protected function requestProcess($request)
    {
        if (isset($request['description'])) {
            $this->apicatContentObj->addParagraph($request['description']);
        }

        $method = isset($request['method']) ? strtoupper($request['method']) : 'GET';
        $methodID = isset($this->requestMethods[$method]) ? $this->requestMethods[$method] : 1;

        $query = [];

        if (isset($request['url'])) {
            if (!isset($request['url']['host']) and !isset($request['url']['path'])) {
                $url = $path = '';
            } else {
                $protocol = (isset($request['url']['protocol']) ? $request['url']['protocol'] : 'http') . '://';
                $host = isset($request['url']['host']) ? implode('.', $request['url']['host']) : '';
                $port = isset($request['url']['port']) ? ':' . $request['url']['port'] : '';
                $path = isset($request['url']['path']) ? '/' . implode('/', $request['url']['path']) : '';
                $url = $protocol . $host . $port;
            }

            if (isset($request['url']['query'])) {
                foreach ($request['url']['query'] as $v) {
                    $query[] = [
                        'name' => $v['key'],
                        'type' => is_numeric($v['value']) ? (strpos($v['value'], '.') !== false ? 2 : 1) : 3,
                        'is_must' => true,
                        'default_value' => $v['value'],
                        'description' => isset($v['description']) ? $v['description'] : '',
                        'sub_params' => []
                    ];
                }
            }
        } else {
            $url = $path = '';
        }

        if (isset($request['body'])) {
            $mode = isset($request['body']['mode']) ? strtolower($request['body']['mode']) : 'none';
            $modeID = isset($this->bodyModes[$mode]) ? $this->bodyModes[$mode] : 0;

            if ($modeID < 1) {
                $body = [];
            } elseif ($modeID < 3 and isset($request['body'][$mode])) {
                $body = $this->baseParamsProcess($request['body'][$mode]);
            } elseif ($modeID == 3) {
                if (!isset($request['body']['options']['raw']['language'], $request['body']['raw'])) {
                    $body = [];
                } else {
                    $body = $this->rawParamsProcess($request['body']['options']['raw']['language'], $request['body']['raw']);
                }
            } else {
                $body = $this->baseParamsProcess([
                    [
                        'key' => 'file',
                        'description' => '',
                        'type' => 'file'
                    ]
                ]);
            }
        } else {
            $modeID = 1;
            $body = [];
        }

        $this->apicatContentObj->addHttpApiUrl($url, $path, $methodID, $modeID);

        $header = [];
        if (isset($request['header']) and $request['header']) {
            $header = $this->baseParamsProcess($request['header']);
        }

        if ($header or $body or $query) {
            $this->apicatContentObj->addHeading('请求参数', 3);
            $this->apicatContentObj->addHttpApiRequestParams(
                'Header 请求参数',
                $header,
                'Body 请求参数',
                $body,
                'Query 请求参数',
                $query
            );
        }
    }

    /**
     * 返回内容处理
     *
     * @param array $response 返回内容
     * @return void
     */
    protected function responseProcess($response)
    {
        if (isset($response['code'])) {
            $this->apicatContentObj->addHttpApiResponseStatusCode($response['code']);
        }

        $header = [];
        if (isset($response['header'])) {
            $header = $this->baseParamsProcess($response['header']);
        }

        $body = [];
        if (isset($response['_postman_previewlanguage'])) {
            $body = $this->rawParamsProcess($response['_postman_previewlanguage'], $response['body']);
        }

        if ($header or $body) {
            $this->apicatContentObj->addHeading('返回参数', 3);
            $this->apicatContentObj->addHttpApiResponseParams('返回头部', $header, '返回参数', $body);
        }

        if (isset($response['_postman_previewlanguage'], $response['body'])) {
            $this->apicatContentObj->addCodeBlock($response['_postman_previewlanguage'], $response['body']);
        }
    }

    /**
     * 常用格式参数处理
     *
     * @param array $params 参数列表
     * @return array
     */
    protected function baseParamsProcess($params)
    {
        $this->paramMakerObj->init();

        foreach ($params as $v) {
            if ($name = strstr($v['key'], '[', true)) {
                // array or object
                $result = preg_match_all('/\[(\w*)\]/', $v['key'], $list);
                if (!$result) {
                    continue;
                }

                array_unshift($list[0], $name);
                $name = $list[0];
            } else {
                // 一级参数
                $name = $v['key'];
            }

            $description = isset($v['description']) ? $v['description'] : '';

            if (isset($v['type']) and $v['type'] == 'file') {
                $this->paramMakerObj->addParam($name, 7, true, '', $description);
            } else {
                if (is_numeric($v['value'])) {
                    $type = strpos($v['value'], '.') !== false ? 2 : 1;
                    $this->paramMakerObj->addParam($name, $type, true, $v['value'], $description);
                } else {
                    $this->paramMakerObj->addParam($name, 3, true, $v['value'], $description);
                }
            }
        }

        return $this->paramMakerObj->params();
    }

    /**
     * raw格式参数处理
     *
     * @param string $language 参数类型
     * @param string $params 参数内容
     * @return array
     */
    protected function rawParamsProcess($language, $params)
    {
        if ($language != 'json' and $language != 'xml') {
            return [];
        }

        if ($language == 'xml') {
            $params = simplexml_load_string($params, 'SimpleXMLElement', LIBXML_NOCDATA);
            $params = json_encode($params);
        }

        if (!$params = json_decode($params)) {
            return [];
        }

        $params = $this->parseJsonToArr($params);

        $this->paramMakerObj->init();

        foreach ($params as $v) {
            $this->paramMakerObj->addParam($v['name'], $v['type'], true, $v['value'], '');
        }

        return $this->paramMakerObj->params();
    }

    /**
     * 解析参数为一维数组
     *
     * @param array $json 参数内容
     * @param array $nameArr 父参数名称数组
     * @return array
     */
    protected function parseJsonToArr($json, $nameArr = [])
    {
        $result = [];

        foreach ($json as $k => $v) {
            $name = is_numeric($k) ? '[]' : ($nameArr ? '[' . $k . ']' : $k);

            if (!is_array($v) and !is_object($v)) {
                if (is_numeric($v)) {
                    $type = strpos($v, '.') !== false ? 2 : 1;
                } else {
                    $type = 3;
                }

                $tmp = $nameArr;
                $tmp[] = $name;

                $result[] = [
                    'name' => $tmp,
                    'type' => $type,
                    'value' => $v
                ];
            } else {
                $tmp = $nameArr;
                $tmp[] = $name;

                $r = $this->parseJsonToArr($v, $tmp);
                if ($r) {
                    $result = array_merge($result, $r);
                }
            }
        }

        return $result;
    }
}
