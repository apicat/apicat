<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Illuminate\Support\Facades\Validator;
use App\Repositories\Project\ProjectRepository;
use App\Repositories\ApiDoc\ApiDocRepository;
use App\Repositories\Project\DocShareRepository;
use App\Repositories\Iteration\IterationRepository;

class HomeController extends Controller
{
    public function index()
    {
        return view('index', [
            'title' => 'ApiCat - API协作开发提效软件',
            'keywords' => 'API开发,API文档,API协作,API调试,API测试,API模拟,API Mock,API导入导出,API对接,文档管理,团队协作,项目协作',
            'description' => 'ApiCat是一款API协作开发提效软件，其简化了团队成员间的协作流程，提供了优质的API文档、API Mock、数据文件导入导出等功能，让开发者可以更快更好的完成开发工作。'
        ]);
    }

    public function project($projectID, $docID = null)
    {
        $validator = Validator::make(
            [
                'project_id' => $projectID,
                'doc_id' => $docID
            ],
            [
                'project_id' => 'required|integer|min:1',
                'doc_id' => 'nullable|integer|min:1'
            ]
        );
        if ($validator->fails()) {
            abort(404);
        }

        if (!$project = ProjectRepository::get($projectID)) {
            abort(404);
        }

        $title = $project->title;
        $keywords = $project->title;
        $description = $project->description;

        if ($docID) {
            if (!$doc = ApiDocRepository::getNode($docID)) {
                abort(404);
            }

            $title .= ' - ' . $doc->title;
            $keywords .= ',' . $doc->title;
        }

        return view('index', [
            'title' => $title,
            'keywords' => $keywords,
            'description' => $description
        ]);
    }

    public function iteration($iterationID, $docID = null)
    {
        $validator = Validator::make(
            [
                'iteration_id' => $iterationID,
                'doc_id' => $docID
            ],
            [
                'iteration_id' => 'required|integer|min:1',
                'doc_id' => 'nullable|integer|min:1'
            ]
        );
        if ($validator->fails()) {
            abort(404);
        }

        if (!$iteration = IterationRepository::get($iterationID)) {
            abort(404);
        }

        $project = ProjectRepository::get($iteration->project_id);
        if (!$project) {
            abort(404);
        }

        $title = $project->title . '(' . $iteration->title . ')';
        $keywords = $project->title;
        $description = $project->description;

        if ($docID) {
            if (!$doc = ApiDocRepository::getNode($docID)) {
                abort(404);
            }

            $title .= ' - ' . $doc->title;
            $keywords .= ',' . $doc->title;
        }

        return view('index', [
            'title' => $title,
            'keywords' => $keywords,
            'description' => $description
        ]);
    }

    public function doc($projectID, $docID)
    {
        $validator = Validator::make(
            [
                'project_id' => $projectID,
                'doc_id' => $docID
            ],
            [
                'project_id' => 'required|integer|min:1',
                'doc_id' => 'required|integer|min:1'
            ]
        );
        if ($validator->fails()) {
            abort(404);
        }

        if (!$project = ProjectRepository::get($projectID)) {
            abort(404);
        }

        if (!$doc = ApiDocRepository::getNode($docID, true)) {
            abort(404);
        }

        return view('index', [
            'title' => $project->title . ' - ' . $doc->title,
            'keywords' => $project->title . ',' . $doc->title,
            'description' => $project->description
        ]);
    }

    public function shareDoc(Request $request, $docID)
    {
        $validator = Validator::make(
            ['doc_id' => $docID],
            ['doc_id' => 'required|integer|min:1']
        );
        if ($validator->fails()) {
            abort(404);
        }

        $share = DocShareRepository::getByDocId($docID);
        if (!$share) {
            abort(404);
        }

        $project = ProjectRepository::get($share->project_id);
        if (!$project) {
            abort(404);
        }

        if (!$doc = ApiDocRepository::getNode($share->doc_id)) {
            abort(404);
        }

        if (substr_count($request->path(), '/') == 1) {
            return view('index', [
                'title' => $doc->title,
                'keywords' => $doc->title,
                'description' => ''
            ]);
        }

        return view('index', [
            'title' => 'ApiCat - API协作开发提效软件',
            'keywords' => 'API开发,API文档,API协作,API调试,API测试,API模拟,API Mock,API导入导出,API对接,文档管理,团队协作,项目协作',
            'description' => 'ApiCat是一款API协作开发提效软件，其简化了团队成员间的协作流程，提供了优质的API文档、API Mock、数据文件导入导出等功能，让开发者可以更快更好的完成开发工作。'
        ]);
    }
}
