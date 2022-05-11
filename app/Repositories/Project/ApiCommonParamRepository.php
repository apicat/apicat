<?php

namespace App\Repositories\Project;

use stdClass;
use App\Models\ApiCommonParam;

class ApiCommonParamRepository
{
    /**
     * 添加默认的公共参数数据
     *
     * @param int $projectID 项目id
     * @return void
     */
    public static function defaultData($projectID)
    {
        $data = [
            ['project_id' => $projectID, 'name' => 'username', 'type' => 3, 'is_must' => 1, 'default_value' => '', 'description' => '用户名'],
            ['project_id' => $projectID, 'name' => 'password', 'type' => 3, 'is_must' => 1, 'default_value' => '', 'description' => '密码'],
            ['project_id' => $projectID, 'name' => 'email', 'type' => 3, 'is_must' => 0, 'default_value' => '', 'description' => '邮箱'],
            ['project_id' => $projectID, 'name' => 'mobile', 'type' => 3, 'is_must' => 0, 'default_value' => '', 'description' => '手机'],
            ['project_id' => $projectID, 'name' => 'page', 'type' => 1, 'is_must' => 0, 'default_value' => '1', 'description' => '页码'],
            ['project_id' => $projectID, 'name' => 'token', 'type' => 3, 'is_must' => 1, 'default_value' => '', 'description' => '访问令牌']
        ];

        foreach ($data as $v) {
            self::create($v);
        }
    }

    /**
     * 创建常用参数
     *
     * @param array $data
     * @return ApiCommonParam
     */
    public static function create($data)
    {
        return ApiCommonParam::create([
            'project_id' => $data['project_id'],
            'name' => $data['name'],
            'type' => $data['type'],
            'is_must' => (int)$data['is_must'],
            'default_value' => $data['default_value'],
            'description' => $data['description']
        ]);
    }

    /**
     * 检查参数名称是否存在
     *
     * @param int $projectID 项目id
     * @param string $name 参数名称
     * @return int
     */
    public static function nameExist($projectID, $name)
    {
        return ApiCommonParam::where([
            ['project_id', $projectID],
            ['name', $name]
        ])->value('id');
    }

    /**
     * 获取常用参数列表
     *
     * @param int $projectID 项目id
     * @return array
     */
    public static function list($projectID)
    {
        $result = ['list' => [], 'map' => new stdClass];

        $records = ApiCommonParam::where('project_id', $projectID)->get();
        if ($records->count() > 0) {
            $paramTypes = ['', 'Int', 'Float', 'String', 'Array', 'Object', 'Boolean', 'File'];

            foreach ($records as $record) {
                $param = [
                    'id' => $record->id,
                    'name' => $record->name,
                    'type' => $record->type,
                    'type_name' => $paramTypes[$record->type],
                    'is_must' => $record->is_must,
                    'default_value' => $record->default_value ?: '',
                    'description' => $record->description ?: ''
                ];

                $paramName = $record->name;
                $result['list'][] = $paramName;
                $result['map']->$paramName = $param;
            }
        }

        return $result;
    }

    /**
     * 获取一条公共参数信息
     *
     * @param int $paramID 参数id
     * @return ApiCommonParam
     */
    public static function get($paramID)
    {
        return ApiCommonParam::find($paramID);
    }

    /**
     * 修改参数内容
     *
     * @param int $paramID 参数id
     * @param array $data 参数信息
     * @return boolean
     */
    public static function update($paramID, $data)
    {
        return (bool)ApiCommonParam::where('id', $paramID)->update([
            'name' => $data['name'],
            'type' => $data['type'],
            'is_must' => (int)$data['is_must'],
            'default_value' => $data['default_value'],
            'description' => $data['description']
        ]);
    }

    /**
     * 删除参数
     *
     * @param int $projectID 项目id
     * @param int $paramID 参数id
     * @return boolean
     */
    public static function remove($projectID, $paramID)
    {
        $param = ApiCommonParam::find($paramID);
        if (!$param or $param->project_id != $projectID) {
            return false;
        }

        return (bool)$param->delete();
    }
}