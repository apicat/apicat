<?php

namespace App\Repositories\Project;

use App\Models\ApiCommonUrl;

class ApiCommonUrlRepository
{
    /**
     * 常用url列表
     * @param int $projectID 项目id
     * @return array
     */
    public static function list($projectID)
    {
        $result = [];

        $records = ApiCommonUrl::where('project_id', $projectID)->get();
        if ($records->count() > 0) {
            foreach ($records as $record) {
                $result[] = [
                    'id' => $record->id,
                    'url' => $record->url
                ];
            }
        }
        return $result;
    }

    /**
     * 添加常用url记录
     * @param int $projectID 项目id
     * @param int $url url
     * @return void
     */
    public static function add($projectID, $url)
    {
        $record = ApiCommonUrl::where([
            ['project_id', $projectID],
            ['url', $url]
        ])->first();

        if (!$record) {
            ApiCommonUrl::create([
                'project_id' => $projectID,
                'url' => $url
            ]);
        }
    }

    /**
     * 删除常用url记录
     * @param int $projectID 项目id
     * @param int $urlID url记录id
     * @return boolean
     */
    public static function remove($projectID, $urlID)
    {
        $record = ApiCommonUrl::find($urlID);
        if (!$record or $record->project_id != $projectID) {
            return false;
        }

        return (bool)$record->delete();
    }
}
