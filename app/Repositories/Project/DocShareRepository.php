<?php

namespace App\Repositories\Project;

use App\Models\DocShare;
use Illuminate\Support\Str;

class DocShareRepository
{
    /**
     * 创建文档分享记录
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @param int $docID 文档id
     * @return DocShare
     */
    public static function create($projectID, $userID, $docID)
    {
        $record = DocShare::where([
            ['project_id', $projectID],
            ['user_id', $userID],
            ['doc_id', $docID]
        ])->first();

        if ($record) {
            $record->secret_key = Str::random(6);
            $record->save();
            return $record;
        }

        return DocShare::create([
            'project_id' => $projectID,
            'user_id' => $userID,
            'doc_id' => $docID,
            'secret_key' => Str::random(6)
        ]);
    }

    /**
     * 获取文档分享详情
     * @param int $projectID 项目id
     * @param int $userID 成员id
     * @param int $docID 文档id
     * @return DocShare
     */
    public static function get($projectID, $userID, $docID)
    {
        return DocShare::where([
            ['project_id', $projectID],
            ['user_id', $userID],
            ['doc_id', $docID]
        ])->first();
    }

    /**
     * 通过文档id获取分享信息
     * @param int $id 分享记录id
     * @return DocShare
     */
    public static function getByDocId($id)
    {
        return DocShare::where('doc_id', $id)->first();
    }

    /**
     * 修改分享的访问密码
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @param int $docID 文档id
     * @return boolean|string 成功: 秘钥  失败: false
     */
    public static function changeSharePassword($projectID, $userID, $docID)
    {
        $secretKey = Str::random(6);

        $res = DocShare::where([
            ['project_id', $projectID],
            ['user_id', $userID],
            ['doc_id', $docID]
        ])->update(['secret_key' => $secretKey]);

        if (!$res) {
            return false;
        }
        return $secretKey;
    }

    /**
     * 删除文档分享记录
     * @param int $projectID 项目id
     * @param int $userID 用户id
     * @param int $docID 文档id
     * @return boolean
     */
    public static function remove($projectID, $userID, $docID)
    {
        return (bool)DocShare::where([
            'project_id' => $projectID,
            'user_id' => $userID,
            'doc_id' => $docID
        ])->delete();
    }
}
