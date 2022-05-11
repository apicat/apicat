<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    /**
     * Run the migrations.
     *
     * @return void
     */
    public function up()
    {
        Schema::create('api_docs', function (Blueprint $table) {
            $table->collation = 'utf8mb4_general_ci';

            $table->id();
            $table->bigInteger('project_id')->index()->comment('项目id');
            $table->bigInteger('parent_id')->comment('父id');
            $table->string('title')->comment('标题');
            $table->tinyInteger('type')->default(0)->comment('类型:0目录,1文档');
            $table->integer('display_order')->default(1)->comment('显示顺序');
            $table->mediumText('content')->nullable()->comment('文档内容');
            $table->bigInteger('created_user_id')->comment('创建者用户id');
            $table->bigInteger('updated_user_id')->comment('上次修改用户id');
            $table->timestamps();

            $table->softDeletes();
        });
    }

    /**
     * Reverse the migrations.
     *
     * @return void
     */
    public function down()
    {
        Schema::dropIfExists('api_docs');
    }
};
