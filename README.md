游戏介绍：

简介：
该游戏是一个抢牌游戏，玩家需要在限定的回合内尽可能多的凑齐能够组成总和为12的卡牌堆，凑齐即可加分，其中游戏中会有4种特殊卡，比如炸弹卡（可以炸掉其他玩家卡堆里的一张卡），万能卡（任意选一张数字类型卡加入自己卡堆），交换卡（可以用自己一张卡与其他玩家交换），修改卡（修改一张自己或者他人的数字类型卡）。卡牌使用方法，玩家可以在游戏中探索。全部回合结束后，根据玩家分数总和计算排名

要求：
玩家卡牌堆的牌会被删除，满足其一：

1.当每一次取到的卡牌放入卡牌堆总和恰好为12时，分数+10，清空卡堆

2.当取到的卡牌放入卡牌堆，总和大于12，此时不加分数，清空卡堆

3.当玩家的卡堆大于6张，则会先删除之前抢的卡，直到卡牌不超过6张

游戏已部署，访问:
8.134.163.22:8001

已搭建分布式系统：

对于配置上，game_web（可开多实例运行）应该部署在同一服务器上（服务器应该存有nginx的图片服务）
user_srv,game_srv，user_web均可部署其他服务器，开多个实例运行
